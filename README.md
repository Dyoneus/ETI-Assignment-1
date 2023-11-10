# ETI-Assignment-1
Assignment 1 project for ETI module

1. Planning of Microservices
Given the requirements of the project, I can split the 2 Microservices into:
- User Service: Manages user accounts, profiles, and authentication.
- Trip Service: Handles the publishing, updating, searching, and booking of trips.

2. Defining Data Structure
For each microservice, define the necessary data structures. For example:
- User Service: User, Profile, CarOwnerProfile entities.
- Trip Service: Trip, SeatReservation entities.



To Do:
BASIC REQUIREMENTS
- Your application must implement the following,
    - Minimum 2 microservices using Go
    - Persistent storage of information using database. e.g. MySQL	
- The data structures, domain, and REST APIs chosen should be appropriate for each microservice.
- Implement a console application that simulates the car-pooling platform frontend. The console application must call your microservices to perform the necessary actions.

BONUS MARKS
- Implement a web frontend that calls your microservices, instead of a console application. You can implement the frontend with any language and design of your choice.

DELIVERABLES
- Submit the following via GitHub (add your tutor as collaborator),
- All source codes of your assignment
- SQL script for setting up your database
- A readme.md indicating
    - Design consideration of your microservices
    - Architecture diagram
    - Instructions for setting up and running your microservices
- Submit a 15-minute presentation of the design and implementation of your assignment. Record your presentation and submit via Brightspace.

