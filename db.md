**Transactions** and **Postings** to mirror hledger's terminology

### transactions Table (plural because "TRANSACTION" is a keyword)
|Column|Type|Description|
|--|:-:|--|
|id|INTEGER PRIMARY KEY|see sqlite datatypes documentation, autoincrementing rowid alias|
|transaction_date|TEXT|see sqlite datetime functions documentation, storing using text|
|last_updated|TEXT|datetime of when the entry was last updated (or inserted)|
|description|TEXT|name/description of transaction|

### posting Table
|Column|Type|Description|
|--|:-:|--|
|id|INTEGER PRIMARY KEY|id of postings|
|transaction_id|INTEGER|foreign key constraint referencing transactions.id|
|account|TEXT|affected account name|
|value|REAL|amount change|
|currency|TEXT|unit of currency|
|comment|TEXT|space for additional comment|

### subscription Table
|Column|Type|Description|
|--|:-:|--|
|id|INTEGER PRIMARY KEY|subscription id|
|description|TEXT|subscription name/description that is used for transaction|
|last_updated|TEXT|date that last transaction was inserted|
|frequency|TEXT|daily, weekly, biweekly, monthly, yearly|

### subscription_template Table
|Column|Type|Description|
|--|:-:|--|
|subscription_id|INTEGER|foreign key constraint referencing subscriptions.id|
|account|TEXT|affected account name|
|value|REAL|amount change|
|currency|TEXT|unit of currency|
|comment|TEXT|space for additional comment|

