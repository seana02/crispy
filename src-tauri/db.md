**Transactions** and **Postings** to mirror hledger's terminology

### crispy Table (metadata)
|Column|Type|Description|
|--|:-:|--|
|key|TEXT|key-value pair key|
|value|TEXT|key-value pair value|

#### Metadata Keys
|Key|Description|
|--|--|
|Version|database structure version|

### transactions Table
|Column|Type|Description|
|--|:-:|--|
|id|INTEGER PRIMARY KEY|see sqlite datatypes documentation, autoincrementing rowid alias|
|transaction_date|TEXT|see sqlite datetime functions documentation, storing using text|
|description|TEXT|name/description of transaction|

### postings Table
|Column|Type|Description|
|--|:-:|--|
|id|INTEGER PRIMARY KEY|id of postings|
|transaction_id|INTEGER|foreign key constraint referencing transactions.id|
|account|TEXT NOT NULL|affected account name|
|value|INTEGER NOT NULL|amount change, multiplied by 10^8 for fixed-point subunit representation (8 digits after the decimal point)|
|currency|TEXT NOT NULL DEFAULT 'USD'|unit of currency|
|comment|TEXT|space for additional comment|

### subscriptions Table
|Column|Type|Description|
|--|:-:|--|
|id|INTEGER PRIMARY KEY|subscription id|
|description|TEXT|subscription name/description that is used for transaction|
|last_updated|TEXT DEFAULT CURRENT_DATE|date that last transaction was inserted|
|frequency|TEXT|daily, weekly, biweekly, monthly, yearly|

### subscription_templates Table
|Column|Type|Description|
|--|:-:|--|
|id|INTEGER PRIMARY KEY|id of template posting|
|subscription_id|INTEGER|foreign key constraint referencing subscriptions.id|
|account|TEXT NOT NULL|affected account name|
|value|INTEGER NOT NULL|amount change|
|currency|TEXT NOT NULL|unit of currency|
|comment|TEXT|space for additional comment|

