# groupie-tracker

## Description

Groupie Tracker is a web-based application that allows users to explore and view detailed information 
about various artists and bands, including their members, activity years, first album release dates, 
concert locations, and upcoming or past concert dates. 

The project uses data from an external API, [Groupie Tracker API](https://groupietrackers.herokuapp.com/api), 
to populate the site with real-time information.

The application provides an interactive and user-friendly interface with data 
visualizations (with cards, and lists) to make it easy for users to navigate through the artists’ details 
and event information. 

Additionally, the application is event-driven, meaning certain actions trigger requests to the server, 
ensuring a dynamic user experience that reflects live data.

## Authors

- **Nicholas Ajwang** - [najwang](https://learn.zone01kisumu.ke/git/najwang)
- **Rabin Odhiambo** - [rotieno](https://learn.zone01kisumu.ke/git/rotieno)

## Usage

### How to Run

To run the Groupie Tracker application, follow these steps:

1. **Clone the Repository**:
    ```shell
    git clone https://learn.zone01kisumu.ke/git/najwang/groupie-tracker.git
    cd groupie-tracker
    ```
2. **Run the Backend**:
    - Navigate to the backend directory and start the server:
      ```shell
      go run main.go
      ```
    - Optionally, you could specify the Port where the application will listen on (by default, the application, listens on port 8080):
      ```shell
      go run main.go -P 9090
      ```
      In the above example, the server starts, listening on port 9090
    
    - If the platform doesn't automatically open on your browser try doing it manually. Open the URL broadcast by the server, in your browser and explore the artists’ information and event data.

## Deployment

A deployed version of the server can be found at: http://groupie-tracker-1.devhive.buzz

## Implementation Details

### Algorithm and Design

1. **Data Fetching and Structuring**:
    - The backend, written in Go, fetches data from the Groupie Tracker API, handling endpoints for `artists`, `locations`, `dates`, and `relation`. Each endpoint provides specific details about artists and their events.
    - The API response is parsed as JSON and structured into Go data types (structs) that align with the API data. This structured data is then exposed to the frontend via HTTP endpoints.

2. **Frontend Data Visualization**:
    - The frontend presents the data with user-friendly visualizations, such as:
        - **Cards** for displaying artist profiles (name, image, first album, members).
        - **Lists** for concert locations and dates.

3. **Event Handling and Client-Server Interaction**:
    - The application is designed to handle client-side events, such as user clicks or filter inputs, that trigger requests to the backend server.
    - When an event (e.g., viewing an artist’s concert history) is triggered, the frontend makes a request to the backend, which processes the request and sends the corresponding data back.

This project demonstrates key concepts in client-server communication, JSON handling, data visualization, and event-driven programming.
