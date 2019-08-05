import { Injectable, EventEmitter } from "@angular/core";

@Injectable({
  providedIn: "root"
})
export class EventService {
  private listener: EventEmitter<any>;
  private retryTimer: number;

  constructor() {
    this.listener = new EventEmitter();
    this.openWebsocket();
  }

  openWebsocket() {
    const endpoint = "ws://" + window.location.host + "/websocket";
    let ws = new WebSocket(endpoint);

    ws.onopen = () => {
      clearInterval(this.retryTimer);
      this.retryTimer = 0;
    };

    ws.onmessage = event => {
      console.log("Emmitting event: " + event);
      this.listener.emit(event);
    };

    ws.onclose = () => {
      ws = null;
      if (!this.retryTimer) {
        this.retryTimer = window.setInterval(() => {
          this.openWebsocket();
        }, 1000);
      }
    };
  }

  getEventListener() {
    return this.listener;
  }
}
