import { Component } from "@angular/core";
import { EventService } from "../services/event.service";
import { MatDialog } from "@angular/material";
import { SuccessDialogComponent } from "../dialogs/success/success-dialog.component";
import { ErrorDialogComponent } from "../dialogs/error/error-dialog.component";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.scss"]
})
export class AppComponent {
  constructor(private events: EventService, private dialog: MatDialog) {
    events.getEventListener().subscribe(event => {
      const data = JSON.parse(event.data);

      const sDialogs = this.dialog.openDialogs.filter(d => {
        return d.componentInstance instanceof SuccessDialogComponent;
      });
      const eDialogs = this.dialog.openDialogs.filter(d => {
        return d.componentInstance instanceof ErrorDialogComponent;
      });

      switch (data.key) {
        case "login":
          if (sDialogs.length === 0) {
            this.dialog.open(SuccessDialogComponent, {
              width: "75vw",
              data: {
                msg: data.data.Name + " successfully checked in!",
                firstName: data.data.FirstName
              }
            });
          }
          break;

        case "login-error":
          if (eDialogs.length === 0) {
            this.dialog.open(ErrorDialogComponent, {
              width: "75vw",
              data: {
                msg: data.data
              }
            });
          }
          break;

        case "card-read-error":
          if (eDialogs.length === 0) {
            this.dialog.open(ErrorDialogComponent, {
              width: "75vw",
              data: {
                msg: "Unable to read card, please try your ID Card again."
              }
            });
          }
          break;
      }
    });
  }
}
