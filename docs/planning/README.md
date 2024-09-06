# Minichat Planning Documentation

## Scope

The goal of our MiniChat project is to create a streamlined and user-friendly
chat application that facilitates real-time communication through channels and direct messages.
We aim to prioritize a clean and intuitive interface while incorporating essential features like
message editing, deletion, and file uploads.
Our focus is on building a solid foundation for efficient and enjoyable communication.

## Resources Required:

- Hosted Server:
  A hosted server is necessary to ensure the app is accessible to users at all times.
  Options include cloud hosting providers like Amazon Web Services (AWS), Google Cloud Platform (GCP) or Hetzner.

- Auth0 License:
  Auth0 provides a robust authentication and authorization platform,
  simplifying the process of user login and securing access to the app's features.
  An Auth0 license is required to utilize their services.

- Domain:
  A custom domain name enhances the app's professional appearance
  and makes it easier for users to remember and access.

# Agile Kanban Board

| Task Name                         | Assignee                        | Estimate (hours) | Size | Priority | Start Date   | End Date     | Milestone           |
| --------------------------------- | ------------------------------- | ---------------- | ---- | -------- | ------------ | ------------ | ------------------- |
| implement chat feature (channels) | CaptainException and MeroFuruya | 3                | M    | P0       |              | Oct 14, 2024 | Develop Application |
| implement event polling/ws        | MeroFuruya                      | 3                | M    | P0       |              | Oct 21, 2024 | Develop Application |
| create navigation sidebar         | CaptainException                | 1                | S    | PO       |              | Oct 7, 2024  | Develop Application |
| create main chat view             | MeroFuruya                      | 1                | S    | PO       |              | Oct 14, 2024 | Develop Application |
| direct messages                   | CaptainException                | 2                | M    | P1       |              | Oct 14, 2024 | Develop Application |
| message markdown support          | CaptainException                | 3                | M    | P2       |              | Nov 4, 2024  | Develop Application |
| file upload for messages          | CaptainException and MeroFuruya | 2                | M    | P1       |              | Nov 11, 2024 | Develop Application |
| file specific "render"            | CaptainException                | 2                | M    | P2       |              | Nov 11, 2024 | Develop Application |
| implement user/channel mentions   | MeroFuruya                      | 1                | S    | P2       |              | Nov 4, 2024  | Develop Application |
| profile customization/settings    | CaptainException                | 3                | L    | P2       |              | Nov 18, 2024 | Develop Application |
| profile picture                   | CaptainException                | 1                | M    | P2       |              | Nov 18, 2024 | Develop Application |
| implement message deletion        | CaptainException                | 0.5              | S    | P1       |              | Oct 28, 2024 | Develop Application |
| Implement message editing         | CaptainException                | 0.5              | S    | P1       |              | Oct 28, 2024 | Develop Application |
| Implement online status           | MeroFuruya                      | 3                | L    | P2       |              | Nov 25, 2024 | Develop Application |
| deploy to server                  | MeroFuruya                      | 2                | M    | PO       |              | Dec 2, 2024  | Develop Application |
| planning-docs                     | CaptainException and MeroFuruya | 4                | L    | P0       | Aug 30, 2024 | Sep 5, 2024  | Planning            |
| create database structure         | CaptainException                | 2                | M    | PO       |              | Sep 30, 2024 | Develop Application |
| set a project goal                | CaptainException and MeroFuruya | 1                | M    | PO       | Aug 29, 2024 | Sep 4, 2024  | Planning            |
| techstack planning                | CaptainException and MeroFuruya | 1                | S    | PO       | Aug 29, 2024 | Sep 2, 2024  | Planning            |
| create tickets                    | CaptainException and MeroFuruya | 1                | S    | PO       | Aug 30, 2024 | Sep 4, 2024  | Planning            |
| setup autho authentication        | MeroFuruya                      | 1                | M    | PO       |              | Sep 23, 2024 | Develop Application |
| pipelines                         | MeroFuruya                      | 1                | M    | P1       | Sep 9, 2024  | Sep 9, 2024  | Develop Application |
| setup eslint-stylistic            | MeroFuruya                      | 0.5              | S    | P1       | Sep 9, 2024  | Sep 9, 2024  | Develop Application |
| Implement users                   | MeroFuruya                      | 2                | M    | P0       |              | Sep 30, 2024 | Develop Application |
| setup nuxt project                | CaptainException and MeroFuruya | 0.5              | XS   | P0       |              | Sep 9, 2024  | Develop Application |
| setup minio                       | MeroFuruya                      | 1.5              | M    | P1       |              | Sep 9, 2024  | Develop Application |
| setup postgres                    | MeroFuruya                      | 1                | S    | P0       |              | Sep 9, 2024  | Develop Application |
| setup database structure          | MeroFuruya                      | 2                | M    | P0       |              | Sep 30, 2024 | Develop Application |
| setup redis                       | MeroFuruya                      | 1                | S    | P2       |              | Nov 25, 2024 | Develop Application |
| setup docker-compose              | MeroFuruya                      | 2                | S    | P0       |              | Dec 2, 2024  | Develop Application |

[GitHub Projects Board](https://github.com/orgs/nineteenseventy/projects/2/views/1)