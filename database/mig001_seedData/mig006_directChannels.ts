import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    insert into minichat.channels_direct (id, user1_id, user2_id) values
      (
        '5c29101e-11a3-4903-84d5-8be262fbeed3'::uuid,
        '038678a5-b6f1-45dc-b1da-9f3837f4cdc8'::uuid,
        '2d0a4682-a7a3-4461-98bc-4b403a94f000'::uuid
      )
  `;
}
