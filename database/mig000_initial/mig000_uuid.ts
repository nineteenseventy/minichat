import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    create extension if not exists "uuid-ossp"
  `;
}
