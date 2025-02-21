### Documentation for the `API` Handler

---

The `API` handler is responsible for handling POST requests to filter a list of artists based on complex filter criteria, which include creation date, first album date, number of band members, and concert locations. Users submit their filter criteria in JSON format, and the server responds with a list of artists that match the criteria.

---

#### **Endpoint**

`POST /api/filter`

---

#### **Request Format**

The request body must be a JSON object that conforms to the structure of the `APIRequestData` type. Below is the structure and explanation of the fields:

```json
{
  "creation_date": {
    "from": 1995,
    "to": 2000,
    "in": [1997, 1998],
    "type": "range"
  },
  "first_album_date": {
    "from": "1990-01-01",
    "to": "1992-12-31",
    "in": ["1991-09-17", "1992-02-10"],
    "type": "range"
  },
  "locations_of_concerts": {
    "in": ["Texas, USA", "Washington, USA"]
  },
  "number_of_members": {
    "from": 4,
    "to": 6,
    "in": [5],
    "type": "or"
  },
  "combinator": "and"
}
```

---

### **Request Fields**

#### **creation_date**
Filter based on the creation date of the artist/band.

- **`from`**: (int) Start of the inclusive date range.
- **`to`**: (int) End of the inclusive date range.
- **`in`**: (array of ints) Specific years to match. Ignored if `type` is `"range"`.
- **`type`**: (string) Determines the filtering logic:
    - `"range"`: Match if creation year is within `[from, to]`.
    - `"in"`: Match if creation year is in the `in` array.
    - `"or"`: Match if either the `range` or `in` filters apply.

#### **first_album_date**
Filter based on the release date of the artist/band's first album.

- **`from`**: (string) Start of the inclusive date range in `YYYY-MM-DD` format.
- **`to`**: (string) End of the inclusive date range in `YYYY-MM-DD` format.
- **`in`**: (array of strings) Specific dates to match. Ignored if `type` is `"range"`.
- **`type`**: (string) Determines the filtering logic:
    - `"range"`: Match if the release date is within `[from, to]`.
    - `"in"`: Match if the release date is in the `in` array.
    - `"or"`: Match if either the `range` or `in` filters apply.

#### **locations_of_concerts**
Filter based on the locations where the artist/band held concerts.

- **`in`** (array of strings): A list of location strings to match.

#### **number_of_members**
Filter based on the number of band members.

- **`from`**: (int) Minimum number of band members (inclusive).
- **`to`**: (int) Maximum number of band members (inclusive).
- **`in`**: (array of ints) Specific counts to match. Ignored if `type` is `"range"`.
- **`type`**: (string) Determines the filtering logic:
    - `"range"`: Match based on the range `[from, to]`.
    - `"in"`: Match if the count is in the `in` array.
    - `"or"`: Match if either the `range` or `in` filters apply.

#### **combinator**
A string that determines the boolean logic across multiple filters. Allowed values:
- `"and"`: All filter conditions must be satisfied.
- `"or"`: At least one filter condition must be satisfied.

---

### **Response Format**

The response is a JSON object with the following structure:

```json
{
  "status": 200,
  "artists": [
    {
      "id": 27,
      "image": "https://example.com/images/bobbyMcFerrin.jpeg",
      "name": "Bobby McFerrins",
      "members": [
        "Bobby McFerrin"
      ],
      "creationDate": 1977,
      "firstAlbum": "1982-09-01",
      "locations": "https://example.com/locations/27",
      "concertDates": "https://example.com/dates/27",
      "relations": "https://example.com/relation/27"
    }
  ]
}
```

---

### **Response Fields**

#### **status**
- HTTP status code of the response (e.g., `200` for success).

#### **artists**
- An array of objects containing details of the artists that matched the filter criteria.

  Each artist object has the following fields:
    - **`id`**: (int) A unique identifier for the artist.
    - **`image`**: (string) URL to the artist's image.
    - **`name`**: (string) The name of the artist/band.
    - **`members`**: (array of strings) Names of the band members.
    - **`creationDate`**: (int) The year the band was formed.
    - **`firstAlbum`**: (string) The release date (format: `YYYY-MM-DD`) of the artist's first album.
    - **`locations`**: (string) API URL with the artist's concert locations.
    - **`concertDates`**: (string) API URL with the artist's concert dates.
    - **`relations`**: (string) API URL with additional artist data relations.

---

### **Examples**

#### Example 1: Filter Artists Created Between 1995 and 2000

**Request**:

```json
POST /api/filter
Content-Type: application/json
{
  "creation_date": {
    "from": 1995,
    "to": 2000,
    "type": "range"
  },
  "combinator": ""
}
```

**Response**:

```json
{
  "status": 200,
  "artists": [
    {
      "id": 30,
      "image": "https://example.com/images/linkinPark.jpeg",
      "name": "Linkin Park",
      "members": ["Chester Bennington", "Mike Shinoda", "Joe Hahn", "Dave Farrell", "Brad Delson", "Rob Bourdon"],
      "creationDate": 1996,
      "firstAlbum": "2000-10-24",
      "locations": "https://example.com/locations/30",
      "concertDates": "https://example.com/dates/30",
      "relations": "https://example.com/relation/30"
    }
  ]
}
```

---

#### Example 2: Filter Artists with Concerts in "Texas, USA"

**Request**:

```json
POST /api/filter
Content-Type: application/json
{
  "locations_of_concerts": {
    "in": ["Texas, USA"]
  },
  "combinator": ""
}
```

**Response**:

```json
{
  "status": 200,
  "artists": [
    {
      "id": 12,
      "image": "https://example.com/images/eminem.jpeg",
      "name": "Eminem",
      "members": ["Marshall Bruce Mathers"],
      "creationDate": 1996,
      "firstAlbum": "1999-10-04",
      "locations": "https://example.com/locations/12",
      "concertDates": "https://example.com/dates/12",
      "relations": "https://example.com/relation/12"
    }
  ]
}
```

---

### **HTTP Status Codes**

- **200 OK**: Success; the `artists` field contains the results.
- **400 Bad Request**: Invalid or malformed request payload.

---

### **Behavior**

The `API` handler processes the filters using the following logic:

1. Validate the request JSON structure. If invalid, respond with a `400 Bad Request` status.
2. Query a predefined list of artists using the filter conditions specified in the request.
3. Combine multiple filters using the specified `combinator` (default is `"or"` if omitted).
4. Return a list of artists that match the filter criteria in the response. If no artists match, return an empty `artists` array.
