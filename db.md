**Transactions** and **Postings** to mirror hledger's terminology

### Transactions Table
|Column|Type|Description|
|--|:-:|--|
|id|INTEGER PRIMARY KEY|see sqlite datatypes documentation, autoincrementing rowid alias|
|transaction_date|TEXT|see sqlite datetime functions documentation, storing using text|
|update_date|INTEGER|datetime of when the entry was entered, in Unix time|
|status|TEXT|complete, pending, scheduled, etc|
|description|TEXT|name/description of transaction|

### Postings Table
|Column|Type|Description|
|--|:-:|--|
|transaction_id|INTEGER|foreign key constraint referencing transactions.id|
|account|TEXT|affected account name|
|value|REAL|amount change|
|currency|TEXT|unit of currency|
|comment|TEXT|space for additional comment|

### Subscriptions Table
|Column|Type|Description|
|--|:-:|--|
|id|INTEGER PRIMARY KEY|subscription id|
|description|TEXT|subscription name/description that is used for transaction|
|last_updated|TEXT|date that last transaction was inserted|
|frequency|TEXT|daily, weekly, biweekly, monthly, yearly|

### Subscription Template Table
|Column|Type|Description|
|--|:-:|--|
|subscription_id|INTEGER|foreign key constraint referencing subscriptions.id|
|account|TEXT|affected account name|
|value|REAL|amount change|
|currency|TEXT|unit of currency|
|comment|TEXT|space for additional comment|

