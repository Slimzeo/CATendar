export namespace agent {

	export class Availability {
	    provider: string;
	    installed: boolean;
	    path: string;
	    version?: string;
	    error?: string;

	    static createFrom(source: any = {}) {
	        return new Availability(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.provider = source["provider"];
	        this.installed = source["installed"];
	        this.path = source["path"];
	        this.version = source["version"];
	        this.error = source["error"];
	    }
	}

}

export namespace ai_sync {

	export class Settings {
	    provider: string;
	    codexPath: string;
	    claudePath: string;

	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.provider = source["provider"];
	        this.codexPath = source["codexPath"];
	        this.claudePath = source["claudePath"];
	    }
	}
	export class AgentConfiguration {
	    settings: Settings;
	    agents: agent.Availability[];

	    static createFrom(source: any = {}) {
	        return new AgentConfiguration(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.settings = this.convertValues(source["settings"], Settings);
	        this.agents = this.convertValues(source["agents"], agent.Availability);
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

	export class Status {
	    runId?: string;
	    state: string;
	    stage?: string;
	    message?: string;
	    provider?: string;
	    emailCount: number;
	    readEmails: number;
	    unreadableEmails: number;
	    created: number;
	    skipped: number;
	    rejected: number;
	    conflicts: number;
	    // Go type: time
	    startedAt?: any;
	    // Go type: time
	    finishedAt?: any;

	    static createFrom(source: any = {}) {
	        return new Status(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.runId = source["runId"];
	        this.state = source["state"];
	        this.stage = source["stage"];
	        this.message = source["message"];
	        this.provider = source["provider"];
	        this.emailCount = source["emailCount"];
	        this.readEmails = source["readEmails"];
	        this.unreadableEmails = source["unreadableEmails"];
	        this.created = source["created"];
	        this.skipped = source["skipped"];
	        this.rejected = source["rejected"];
	        this.conflicts = source["conflicts"];
	        this.startedAt = this.convertValues(source["startedAt"], null);
	        this.finishedAt = this.convertValues(source["finishedAt"], null);
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

export namespace calendar {

	export class Event {
	    id: number;
	    title: string;
	    start: string;
	    end?: string;
	    color: string;
	    allDay: boolean;
	    description: string;
	    bold: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;

	    static createFrom(source: any = {}) {
	        return new Event(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.start = source["start"];
	        this.end = source["end"];
	        this.color = source["color"];
	        this.allDay = source["allDay"];
	        this.description = source["description"];
	        this.bold = source["bold"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class EventInput {
	    title: string;
	    start: string;
	    end?: string;
	    color: string;
	    allDay: boolean;
	    description: string;
	    bold: boolean;

	    static createFrom(source: any = {}) {
	        return new EventInput(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.start = source["start"];
	        this.end = source["end"];
	        this.color = source["color"];
	        this.allDay = source["allDay"];
	        this.description = source["description"];
	        this.bold = source["bold"];
	    }
	}

}

export namespace email {

	export class Account {
	    id: string;
	    address: string;
	    username: string;
	    imapHost: string;
	    imapPort: number;
	    folder: string;
	    useTLS: boolean;
	    configured: boolean;
	    hasCredential: boolean;

	    static createFrom(source: any = {}) {
	        return new Account(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.address = source["address"];
	        this.username = source["username"];
	        this.imapHost = source["imapHost"];
	        this.imapPort = source["imapPort"];
	        this.folder = source["folder"];
	        this.useTLS = source["useTLS"];
	        this.configured = source["configured"];
	        this.hasCredential = source["hasCredential"];
	    }
	}
	export class AccountInput {
	    address: string;
	    username: string;
	    imapHost: string;
	    imapPort: number;
	    folder: string;
	    useTLS: boolean;
	    secret: string;

	    static createFrom(source: any = {}) {
	        return new AccountInput(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.address = source["address"];
	        this.username = source["username"];
	        this.imapHost = source["imapHost"];
	        this.imapPort = source["imapPort"];
	        this.folder = source["folder"];
	        this.useTLS = source["useTLS"];
	        this.secret = source["secret"];
	    }
	}
	export class Message {
	    id: string;
	    messageId: string;
	    subject: string;
	    sender: string;
	    // Go type: time
	    receivedAt: any;
	    // Go type: time
	    sentAt?: any;
	    sourceUrl: string;
	    text: string;
	    readError?: string;

	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.messageId = source["messageId"];
	        this.subject = source["subject"];
	        this.sender = source["sender"];
	        this.receivedAt = this.convertValues(source["receivedAt"], null);
	        this.sentAt = this.convertValues(source["sentAt"], null);
	        this.sourceUrl = source["sourceUrl"];
	        this.text = source["text"];
	        this.readError = source["readError"];
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

export namespace main {

	export class AISyncSetup {
	    email: email.Account;
	    agent: ai_sync.AgentConfiguration;

	    static createFrom(source: any = {}) {
	        return new AISyncSetup(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.email = this.convertValues(source["email"], email.Account);
	        this.agent = this.convertValues(source["agent"], ai_sync.AgentConfiguration);
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
