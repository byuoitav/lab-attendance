class APIService {
  constructor() {}

  login(byuID) {
    const endpoint = `http://localhost:8243/api/v1/login/${byuID}`;

    console.log("Hitting endpoint: " + endpoint);

    fetch(endpoint, {
      method: "POST"
    })
      .then(response => {
        if (!response.ok) {
          throw new Error("Network response was not ok: " + response.statusText);
        }

        return;
      })
      .then(data => {
        console.log(data);
      })
      .catch(error => {
        console.log("Error:", error);
      });
  }

  getLabName() {
    const endpoint = "http://localhost:8243/api/v1/config";

    return fetch(endpoint)
      .then(response => {
        if (!response.ok) {
          throw new Error("Network response was not ok: " + response.statusText);
        }
        return response.json();
      })
      .then(data => {
        return data.lab_name;
      })
      .catch(error => {
        console.log("Error:", error);
        return null;
      });
  }
}

class EventService {
  constructor() {
    this.listener = new EventTarget();
    this.retryTimer = null;
    this.openWebsocket();
  }

  openWebsocket() {
    const endpoint = "ws://localhost:8243/websocket";
    let ws = new WebSocket(endpoint);

    ws.onopen = () => {
      console.log("WebSocket connected");
      clearInterval(this.retryTimer);
      this.retryTimer = null;
    };

    ws.onmessage = (event) => {
      console.log("Emitting event:", event);
      const customEvent = new CustomEvent("message", { detail: event });
      this.listener.dispatchEvent(customEvent);
    };

    ws.onclose = () => {
      console.log("WebSocket closed, attempting reconnect...");
      ws = null;
      if (!this.retryTimer) {
        this.retryTimer = setInterval(() => {
          this.openWebsocket();
        }, 1000);
      }
    };
  }

  getEventListener() {
    return this.listener;
  }
}

// Example usage:
// const events = new EventService();
// events.getEventListener().addEventListener("message", (event) => {
//   console.log("Received data:", event.detail.data);
// });
