import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    insert into minichat.users (id, idp_id, username, color, bio) values
      (
        '038678a5-b6f1-45dc-b1da-9f3837f4cdc8'::uuid,
        'auth0|66d4ff54bd853b4ba1a9edfd',
        'MeroFuruya',
        '#FF0000',
        'I am a software engineer.'
      ),
      (
        '2d0a4682-a7a3-4461-98bc-4b403a94f000'::uuid,
        'auth0|66d6c778063ba07e03e94d29',
        'CaptainChrom',
        '#00FF00',
        'I am a software engineer.'
      )
  `;
}
