import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    insert into minichat.channels_public (id, title, description) values
      (
        '4b7c7f2a-8d8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        'General',
        'General discussion'
      ),
      (
        'f7d7a4b5-8f8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        'Random',
        'Random discussion'
      ),
      (
        'b7d7a4b5-8f8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        'Tech',
        'Tech discussion'
      ),
      (
        'd7d7a4b5-8f8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        'Music',
        'Music discussion'
      ),
      (
        'e7d7a4b5-8f8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        'Movies',
        'Movies discussion'
      ),
      (
        'c7d7a4b5-8f8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        'Books',
        'Books discussion'
      ),
      (
        'a7d7a4b5-8f8d-4c4f-b4a7-6f5c7b1f5d3e'::uuid,
        'Games',
        'Games discussion'
      )
  `;
}
