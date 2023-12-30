CRUD

- GET
- POST
  - [] port over /api/validate_chirp to the /api/chirps call -> keep the validation
  - [] add unique id param - increment - e.g. 1, 2, 3, 4
    - request:
      {
      "body": "Hello, world!"
      }
    - reponse:
      {
      "id": 1,
      "body": "Hello, world!"
      }
- Saving to the DB
  - [] When updating the database, read the entire thing into memory (unmarshal it into a struct), update the data, and then write the entire thing back to disk (marshal it back into JSON).
  - [] To make sure that multiple requests don't try to write to the database at the same time, use a mutex to lock the database while you're using it.

Std lib pkgs:

- os.ReadFile
- os.ErrNotExist
- os.WriteFile
- sort.Slice
