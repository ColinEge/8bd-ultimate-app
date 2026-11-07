export namespace service {
	
	export class Action {
	    actionCode: string;
	    actionLabel: string;
	
	    static createFrom(source: any = {}) {
	        return new Action(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.actionCode = source["actionCode"];
	        this.actionLabel = source["actionLabel"];
	    }
	}
	export class ActionCategory {
	    Name: string;
	    Actions: Action[];
	
	    static createFrom(source: any = {}) {
	        return new ActionCategory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Actions = this.convertValues(source["Actions"], Action);
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
	export class Device {
	    DeviceID: string;
	    ProductName: string;
	
	    static createFrom(source: any = {}) {
	        return new Device(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DeviceID = source["DeviceID"];
	        this.ProductName = source["ProductName"];
	    }
	}
	export class Key {
	    id: string;
	    physicalKey: number;
	    keyLabel: string;
	    bindingState: string;
	    actionCode: string;
	    actionLabel: string;
	
	    static createFrom(source: any = {}) {
	        return new Key(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.physicalKey = source["physicalKey"];
	        this.keyLabel = source["keyLabel"];
	        this.bindingState = source["bindingState"];
	        this.actionCode = source["actionCode"];
	        this.actionLabel = source["actionLabel"];
	    }
	}
	export class Keyboard {
	    Model: number;
	    DeviceID: string;
	    ProductName: string;
	    HasKeypad: boolean;
	    Keys: Record<string, Key>;
	
	    static createFrom(source: any = {}) {
	        return new Keyboard(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Model = source["Model"];
	        this.DeviceID = source["DeviceID"];
	        this.ProductName = source["ProductName"];
	        this.HasKeypad = source["HasKeypad"];
	        this.Keys = this.convertValues(source["Keys"], Key, true);
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
	export class ProfileBinding {
	    key: number;
	    action: string;
	
	    static createFrom(source: any = {}) {
	        return new ProfileBinding(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.action = source["action"];
	    }
	}
	export class Profile {
	    title: string;
	    keys: ProfileBinding[];
	
	    static createFrom(source: any = {}) {
	        return new Profile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.keys = this.convertValues(source["keys"], ProfileBinding);
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

