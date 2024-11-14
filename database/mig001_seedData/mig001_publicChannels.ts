import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    insert into minichat.channels (id, type) values
      (
        '4b7c7f2a-8d8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        'public'
      ),
      (
        'f7d7a4b5-8f8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        'public'
      ),
      (
        'b7d7a4b5-8f8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        'public'
      ),
      (
        'd7d7a4b5-8f8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        'public'
      ),
      (
        'e7d7a4b5-8f8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        'public'
      ),
      (
        'c7d7a4b5-8f8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        'public'
      ),
      (
        'a7d7a4b5-8f8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        'public'
      )
  `;
}
