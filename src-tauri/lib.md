aggregate/summary information:
- previous month income
- current month income
- current liabilities
- current month spending

filters:
- date range
- account
- text search (transaction description and posting comment)

functions:
get_all_transaction_list                            get a list of transactions with no posting information
get_transaction_by_id                               get complete information about a single transaction
get_transactions_between                            get a list of transactions with no posting information between the given dates
get_transactions_involving_accounts                 get a list of transactions with no posting information involving the given account
get_transactions_with_description                   get a list of transactions with no posting information containing a text string in the description or comments
create_transaction                                  insert transaction+postings by id
delete_transaction                                  delete transaction+postings by id
update_transaction                                  update transaction+postings by id
