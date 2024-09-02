import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
  alter table minichat.users
  rename column profile_picture to picture;
  `;
}
