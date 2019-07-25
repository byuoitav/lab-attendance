import { Component } from "@angular/core";
import { Router } from "@angular/router";
import { APIService } from "../../services/api.service";

@Component({
  selector: "login",
  templateUrl: "./login.component.html",
  styleUrls: ["./login.component.scss"]
})
export class LoginComponent {
  id = "";
  ssCounter = 0;
  ssTimeoutMax = 30;
  ssTimer: any;

  constructor(private router: Router, private api: APIService) {
    this.ssTimer = setInterval(() => {
      this.ssCounter++;
      // console.log("counter", this.ssCounter);

      if (this.ssCounter >= this.ssTimeoutMax) {
        this.ssCounter = 0;
        clearInterval(this.ssTimer);
        this.router.navigate(["/screensaver"]);
      }
    }, 1000);
  }

  addToID(num: string) {
    if (this.id.length < 9) {
      this.id += num;
    }
  }

  delFromID() {
    if (this.id.length > 0) {
      this.id = this.id.slice(0, -1);
    }
  }

  login = async (id: string) => {
    console.log("logging in id", this.id);
    this.ssCounter = 0;

    this.api.login(id);

    this.id = ""; // reset the id
  };
}
