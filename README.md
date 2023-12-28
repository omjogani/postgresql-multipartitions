# PostgreSQL Multi-Partitions
Multi-partition architecture which divides the traffic according to the request type. The server will decide to navigate the request to which Database based on the Request Type. In this example, There are 2 Master Database in the PostgreSQL Database & Golang Server will decide which Database to forward the request to.

### Architecture
![architectural_figure](https://github.com/omjogani/postgresql-multipartitions/blob/master/postgresql-multipartition.png?raw=true "Architectural Figure")

### Working Example

![multipartition-postgresql-demo](https://github.com/omjogani/postgresql-multipartitions/assets/72139914/4cd5f495-a127-4508-a12d-a122414ab643)

### Run Locally
> Make sure that you have Golang installed in your System!
- Rename `.example-env` to `.env` and add PostGre Username, Password & DB Name.
- Run the command:
```
go run github.com/cosmtrek/air
```
- Getting an error in above mentioned command? Run...
```
go get github.com/cosmtrek/air
```
- then run the First command.


### Technical Details

---
- Technology
    - Server: Golang
    - Database: PostgreSQL
    - User Interface: HTML, Tailwind CSS


>If you found this useful, make sure to give it a star 🌟
## Thank You!!
