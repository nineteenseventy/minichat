# Channel

- Do we want to display channel members for public channels? - no
- Do we want to display "channel members" for direct channels? - no
- Extra endpoint for (group) channel members? - i suppose yes
- 

## Channel Members

- Endpoint: `/channels/{channel_id}/members`
- Method: `GET`
- Should only work for group channels that caller is a member of
- Returns: List of users in the channel
