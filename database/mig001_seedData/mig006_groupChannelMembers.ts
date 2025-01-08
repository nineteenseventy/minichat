import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    insert into minichat.channels_members (id, channel_id, user_id) values
      (
        '1cf8524d-7006-4adc-ac67-bd67f3684c81'::uuid,
        '747bb027-6761-4571-bfe6-96ec3d866b5b'::uuid,
        '038678a5-b6f1-45dc-b1da-9f3837f4cdc8'::uuid
      ),
      (
        'afac5fb3-6f6e-40f3-b5ad-34b1288ffa28'::uuid,
        '747bb027-6761-4571-bfe6-96ec3d866b5b'::uuid,
        '2d0a4682-a7a3-4461-98bc-4b403a94f000'::uuid
      ),
      (
        'a0c45ec6-15a4-4669-915a-ce08c388b2c2'::uuid,
        '01cc33cf-8e71-4a08-8bc8-137d91c107ac'::uuid,
        '038678a5-b6f1-45dc-b1da-9f3837f4cdc8'::uuid
      ),
      (
        '2bbed69a-6835-4936-90cc-ec407cf8d9db'::uuid,
        '01cc33cf-8e71-4a08-8bc8-137d91c107ac'::uuid,
        '2d0a4682-a7a3-4461-98bc-4b403a94f000'::uuid
      )
  `;
}
