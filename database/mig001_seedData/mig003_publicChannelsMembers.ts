import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    insert into minichat.channels_members (id, channel_id, user_id) values
      (
        '5506fdb0-4efe-4b53-94a7-5446c6bc90e6'::uuid,
        '4b7c7f2a-8d8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        '038678a5-b6f1-45dc-b1da-9f3837f4cdc8'::uuid
      ),
      (
        '296d59a4-f92e-482f-baa0-16a14604f91d'::uuid,
        '4b7c7f2a-8d8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        '2d0a4682-a7a3-4461-98bc-4b403a94f000'::uuid
      ),
      (
        'ef01ae5b-a9af-46e6-a7e2-6e9264d50226'::uuid,
        'b7d7a4b5-8f8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        '038678a5-b6f1-45dc-b1da-9f3837f4cdc8'::uuid
      ),
      (
        '501bb5ca-2242-4729-b30a-fd0c436f9a89'::uuid,
        'b7d7a4b5-8f8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        '2d0a4682-a7a3-4461-98bc-4b403a94f000'::uuid
      ),
      (
        '2071bc1e-8003-4ea7-a42b-36301f65bddb'::uuid,
        'e7d7a4b5-8f8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        '038678a5-b6f1-45dc-b1da-9f3837f4cdc8'::uuid
      ),
      (
        '18637121-cb80-4d9a-916a-c08a11c16a1a'::uuid,
        'e7d7a4b5-8f8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        '2d0a4682-a7a3-4461-98bc-4b403a94f000'::uuid
      )
  `;
}
