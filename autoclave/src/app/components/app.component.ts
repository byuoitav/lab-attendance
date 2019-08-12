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

      if (data.key === "login") {
        if (data.value === "true") {
          const sDialogs = this.dialog.openDialogs.filter(d => {
            return d.componentInstance instanceof SuccessDialogComponent;
          });

          if (sDialogs.length === 0) {
            this.dialog.open(SuccessDialogComponent, {
              width: "75vw",
              data: {
                msg: data.data + " (" + data.user + ") successfully checked in!"
              }
            });
          }
          console.log(data);
        } else {
          const eDialogs = this.dialog.openDialogs.filter(d => {
            return d.componentInstance instanceof ErrorDialogComponent;
          });

          if (eDialogs.length === 0) {
            this.dialog.open(SuccessDialogComponent, {
              width: "75vw",
              data: {
                msg: data.data
              }
            });
          }
        }
      }

      if (data.key === "card-read-error") {
        const eDialogs = this.dialog.openDialogs.filter(d => {
          return d.componentInstance instanceof ErrorDialogComponent;
        });

        if (eDialogs.length === 0) {
          this.dialog.open(SuccessDialogComponent, {
            width: "75vw",
            data: {
              msg: "Unable to read card, please try your ID Card again."
            }
          });
        }
      }
    });
  }
}
