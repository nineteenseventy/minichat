import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    insert into minichat.messages (id, author_id, channel_id, content) values
      (
        'b25c06cd-57ba-4288-96c2-baf85bc51cca'::uuid,
        '038678a5-b6f1-45dc-b1da-9f3837f4cdc8'::uuid,
        '5c29101e-11a3-4903-84d5-8be262fbeed3'::uuid,
        'Test Message 1 in Direct Channel by MeroFuruya'
      ),
      (
        '2cb0d2f2-331a-41e0-ad75-3160ab6b5e02'::uuid,
        '2d0a4682-a7a3-4461-98bc-4b403a94f000'::uuid,
        '4b7c7f2a-8d8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        'Test Message 2 in Public Channel by CaptainChrom'
      )
  `;
}
