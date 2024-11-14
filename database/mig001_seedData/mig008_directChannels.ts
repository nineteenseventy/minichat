import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    insert into minichat.channels_direct (id) values
      (
        '5c29101e-11a3-4903-84d5-8be262fbeed3'::uuid
      )
  `;
}
