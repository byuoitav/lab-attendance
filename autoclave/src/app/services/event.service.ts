import { Injectable, EventEmitter } from "@angular/core";

@Injectable({
  providedIn: "root"
})
export class EventService {
  private listener: EventEmitter<any>;

  constructor() {
    const endpoint = "ws://" + window.location.host + "/websocket";
    const ws = new WebSocket(endpoint);
    this.listener = new EventEmitter();

    ws.onmessage = event => {
      console.log("Emmitting event: " + event);
      this.listener.emit(event);
    };
  }

  getEventListener() {
    return this.listener;
  }
}
