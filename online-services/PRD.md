# Product Requirements Document (PRD)

## 1. Purpose

The purpose of this PRD is to outline the requirements for the online component of the project, which includes player connectivity, statistics tracking, achievements, matchmaking, and optional features like an admin interface and in-game store. This document will serve as a guide for the development team to ensure all necessary features are implemented according to the project's goals and user needs.

## 2. Objectives

- **Player Connectivity**: Implement a secure and efficient system for players to connect to the game using JSON Web Tokens (JWT) for authentication.
- **Statistics Tracking**: Develop a system to track and store player statistics, such as the number of games won and cubes cleared, using an in-memory map with the option to upgrade to Redis or a database.
- **Achievements**: Create a system to award achievements based on player statistics and milestones.
- **Matchmaking**: Implement a matchmaking system that pairs players of similar skill levels and directs them to the same server.
- **Optional Features**:
  - **Admin Interface**: Develop an interface for administrators to manage achievements, statistics, and view ongoing games.
  - **In-Game Store**: Create a store for players to purchase cosmetic items.

## 3. Key Features

### 3.1 Player Connectivity

- **JWT Authentication**: Use JSON Web Tokens (JWT) to authenticate and recognize players.
- **Secure Connection**: Ensure that the connection process is secure and efficient.
- **Iterative Development**: Focus on quick iterations rather than a perfect solution initially.

### 3.2 Statistics Tracking

- **In-Memory Map**: Use an in-memory map to store statistics initially.
- **Upgrade Options**: Plan for future upgrades to Redis or a database for better performance and scalability.
- **Tracked Metrics**: Include metrics such as the number of games won and cubes cleared.

### 3.3 Achievements

- **Milestone-Based**: Award achievements based on specific milestones, such as winning a certain number of games or clearing a certain number of cubes.
- **Hardcoded Links**: Initially, link achievements to statistics in a hardcoded manner.

### 3.4 Matchmaking

- **Server Registration**: Servers register with the matchmaking system when they start.
- **Player Pairing**: Pair players of similar skill levels and direct them to the same server.
- **Skill Level Definition**: Define skill levels based on relevant criteria.

### 3.5 Optional Features

#### 3.5.1 Admin Interface

- **Achievement Management**: Allow administrators to add and manage achievements.
- **Statistics Management**: Allow administrators to add and manage statistics.
- **Game Monitoring**: Provide a list of ongoing games for administrators to monitor.

#### 3.5.2 In-Game Store

- **Cosmetic Items**: Allow players to purchase cosmetic items to customize their game experience.

## 4. Assumptions

- **Technology Stack**: The system will be developed using Golang.
- **API Design**: The online component will be designed as a REST API.
- **User Base**: Assume a small (<1000) user base with varying skill levels.
- **Security**: Assume that security is a top priority, especially for player authentication and data storage.

## 5. Success Metrics

- **Player Connectivity**: Successful and secure player authentication and connection.
- **Statistics Tracking**: Accurate and efficient tracking of player statistics.
- **Achievements**: Proper awarding of achievements based on milestones.
- **Matchmaking**: Effective pairing of players based on skill levels.
- **Optional Features**: Functional admin interface and in-game store.

## 6. User Scenarios

### 6.1 Player Connectivity

- **Scenario**: A player logs in to the game and is authenticated using JWT.
- **Expected Outcome**: The player is successfully connected and recognized by the system.

### 6.2 Statistics Tracking

- **Scenario**: A player wins a game and clears a certain number of cubes.
- **Expected Outcome**: The system accurately tracks and stores the player's statistics.

### 6.3 Achievements

- **Scenario**: A player reaches a milestone, such as winning 10 games.
- **Expected Outcome**: The system awards the player with the corresponding achievement.

### 6.4 Matchmaking

- **Scenario**: A player requests a game, and the system pairs them with other players of similar skill levels.
- **Expected Outcome**: The player is directed to a server with similarly skilled players.

### 6.5 Optional Features

#### 6.5.1 Admin Interface

- **Scenario**: An administrator adds a new achievement to the system.
- **Expected Outcome**: The achievement is successfully added and available to players.

#### 6.5.2 In-Game Store

- **Scenario**: A player purchases a cosmetic item from the in-game store.
- **Expected Outcome**: The transaction is successful, and the player receives the item.

## 7. Technical Specifications

- **Languages**: Golang
- **API Design**: REST API
- **Authentication**: JSON Web Tokens (JWT)
- **Data Storage**: SQLite database
- **Security**: High priority on secure player authentication and data storage

## 8. Timeline

- **Phase 1**: Implement player connectivity and JWT authentication.
- **Phase 2**: Develop statistics tracking and achievements.
- **Phase 3**: Implement matchmaking system.
- **Phase 4**: Develop optional features (admin interface and in-game store).

## 9. Risks and Challenges

- **Security Risks**: Ensuring secure player authentication and data storage.
- **Matchmaking Accuracy**: Ensuring the matchmaking system accurately pairs players based on skill levels.