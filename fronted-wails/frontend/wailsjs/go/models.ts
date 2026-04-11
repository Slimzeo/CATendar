export namespace models {
	
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

