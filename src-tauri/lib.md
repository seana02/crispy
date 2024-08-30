DB operations module
- CRUD transactions
    - C: 1 Transaction obj + many Posting obj
        - validate zero-sum (separate function for reuse)
        - DB TRIGGER rules for auto postings (see hledger docs)
    - R: by id, search by desc (multiple)
    - U: separate funcs per field, chosen by caller
        - or id and "new Transaction object" with validation check
    - D: by id
