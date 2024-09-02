import fs from 'fs';
import { fileURLToPath } from 'url';
import path from 'path';
import postgres from 'postgres';
import dotenv from 'dotenv';

const PG_MIN_SAFE_INTEGER = -2147483648;

export interface MigrationFunction {
  (sql: postgres.Sql): Promise<void>;
}

interface Migration {
  migrationNumber: number;
  fileNumber: number;
  name: string;
  filename: string;
}

let logNotice = false;
function onNotice(notice: postgres.Notice) {
  if (logNotice) console.log(notice);
}

function filterAndSortMigrationNames(names: string[]) {
  return names
    .reduce<{ name: string; number: number }[]>((acc, name) => {
      const match = /^mig(?<number>[0-9]*)_/.exec(name);
      if (match && match.groups && match.groups.number) {
        acc.push({
          name: name,
          number: Number.parseInt(match.groups.number, 10),
        });
      }
      return acc;
    }, [])
    .sort((a, b) => a.number - b.number);
}

function listMigrationFiles(): Migration[] {
  const allFolders = fs.readdirSync(
    path.join(path.dirname(fileURLToPath(import.meta.url))),
  );
  const folders = filterAndSortMigrationNames(allFolders);

  const migrationFiles: Migration[] = [];

  for (const folder of folders) {
    const allFiles = fs.readdirSync(
      path.join(path.dirname(fileURLToPath(import.meta.url)), folder.name),
    );
    const files = filterAndSortMigrationNames(allFiles);
    migrationFiles.push(
      ...files.map((file) => ({
        migrationNumber: folder.number,
        fileNumber: file.number,
        name: folder.name,
        filename: file.name,
      })),
    );
  }
  return migrationFiles;
}

async function run(migration: Migration, sql: postgres.Sql) {
  const { migrationNumber, fileNumber, name, filename } = migration;
  const module = await import(`./${name}/${filename}`).catch((error) => {
    console.error(error);
    process.exit(1);
  });

  if (!module.default) {
    throw new Error(
      `Migration ${name}/${filename} does not export a default function`,
    );
  }

  const migrationFunction = module.default as MigrationFunction;

  await migrationFunction(sql).catch((error) => {
    throw new Error(`Error running migration ${name}/${filename}: ${error}`);
  });

  await sql`
    insert into migrations (number, file_number, name, filename)
    values (${migrationNumber}, ${fileNumber}, ${name}, ${filename})
  `.catch((error) => {
    throw new Error(
      `Error inserting migration ${name}/${filename} into database: ${error}`,
    );
  });
  console.log(`Ran migration ${name}/${filename}`);
}

async function ensureMigrationsTable(sql: postgres.Sql) {
  await sql`create extension if not exists "uuid-ossp"`.execute();
  await sql`
      create table if not exists migrations (
        id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
        created_at timestamp with time zone not null default now(),
        number integer not null,
        file_number integer not null,
        name text not null,
        filename text not null
      )
    `.execute();
}

async function main() {
  dotenv.config();
  const sql = postgres({
    onnotice: onNotice,
  });
  await ensureMigrationsTable(sql);
  logNotice = true;

  const highestMigration = await sql`
  select Max(number)
  from migrations
  `.then((result) => result[0].max ?? PG_MIN_SAFE_INTEGER);

  const highestMigrationFile = await sql`
    select Max(file_number)
    from migrations
    where number = ${highestMigration}
  `.then((result) => result[0].max ?? PG_MIN_SAFE_INTEGER);

  const migrations = listMigrationFiles().filter(
    (migration) =>
      migration.migrationNumber > highestMigration ||
      (migration.migrationNumber === highestMigration &&
        migration.fileNumber > highestMigrationFile),
  );

  for (const migration of migrations) {
    await run(migration, sql);
  }

  await sql.end();
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
