export namespace domain {
	
	export enum Status {
	    Pending = "Pending",
	    Posted = "Posted",
	    Cleared = "Cleared",
	}
	export enum Type {
	    Asset = "Asset",
	    Liability = "Liability",
	    Revenue = "Revenue",
	    Expense = "Expense",
	    Equity = "Equity",
	}
	export class AccountDTO {
	    id: number;
	    parentID: number;
	    name: string;
	    type: Type;
	    currency: string;
	    description: string;
	    active: boolean;
	    dateCreated: string;
	    dateUpdated: string;
	
	    static createFrom(source: any = {}) {
	        return new AccountDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.parentID = source["parentID"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.currency = source["currency"];
	        this.description = source["description"];
	        this.active = source["active"];
	        this.dateCreated = source["dateCreated"];
	        this.dateUpdated = source["dateUpdated"];
	    }
	}
	export class PostingDTO {
	    id: number;
	    transactionID: number;
	    accountID: number;
	    accountName: string;
	    // Go type: decimal
	    amount: any;
	    currency: string;
	    dateCreated: string;
	    dateUpdated: string;
	
	    static createFrom(source: any = {}) {
	        return new PostingDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.transactionID = source["transactionID"];
	        this.accountID = source["accountID"];
	        this.accountName = source["accountName"];
	        this.amount = this.convertValues(source["amount"], null);
	        this.currency = source["currency"];
	        this.dateCreated = source["dateCreated"];
	        this.dateUpdated = source["dateUpdated"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TransactionDTO {
	    id: number;
	    description: string;
	    date: string;
	    status: Status;
	    referenceID?: number;
	    dateCreated: string;
	    dateUpdated: string;
	    postings: PostingDTO[];
	    tags: string[];
	
	    static createFrom(source: any = {}) {
	        return new TransactionDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.description = source["description"];
	        this.date = source["date"];
	        this.status = source["status"];
	        this.referenceID = source["referenceID"];
	        this.dateCreated = source["dateCreated"];
	        this.dateUpdated = source["dateUpdated"];
	        this.postings = this.convertValues(source["postings"], PostingDTO);
	        this.tags = source["tags"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

