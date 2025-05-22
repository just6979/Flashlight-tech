Go for back-end code, PostgreSQL as the database engine, and React (written in TypeScript and using Vite for builds) for the front-end, but this is only a suggestion.


Back-end: Implement a RESTful API with CRUD operations for managing students.
* API
  * GET /students - Retrieve a list of all students.
  * POST /students - Add a new student.
  * PUT /students/:id - Update an existing student.
  * DELETE /students/:id - Delete a student.
* Data Model
  * id (integer, primary key)
  * name (string)
  * grade (integer)
* Front-end: Create a functional UI that interacts with the backend API. It does not need to be polished or aesthetically pleasing; it just needs to be able to connect with the back-end and perform the following functions:
  * View a list of all students.
  * Add a new student.
   * Update an existing student.
  * Delete a student.
