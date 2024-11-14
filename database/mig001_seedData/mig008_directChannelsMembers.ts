import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    insert into minichat.channels_members (id, channel_id, user_id) values
      (
        '4a4a340c-962b-4b41-ab55-d36a34533a19'::uuid,
        '5c29101e-11a3-4903-84d5-8be262fbeed3'::uuid,
        '038678a5-b6f1-45dc-b1da-9f3837f4cdc8'::uuid
      ),
      (
        '2c102ba3-a654-40e2-aece-6ee7381633a5'::uuid,
        '5c29101e-11a3-4903-84d5-8be262fbeed3'::uuid,
        '2d0a4682-a7a3-4461-98bc-4b403a94f000'::uuid
      )
  `;
}
