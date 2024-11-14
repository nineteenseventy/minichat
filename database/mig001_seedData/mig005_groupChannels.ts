import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    insert into minichat.channels_group (id, title) values
      (
        '747bb027-6761-4571-bfe6-96ec3d866b5b'::uuid,
        'Tech Group'
      ),
      (
        '01cc33cf-8e71-4a08-8bc8-137d91c107ac'::uuid,
        'Games Group'
      )
  `;
}
