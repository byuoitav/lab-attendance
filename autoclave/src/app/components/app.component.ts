import { Component } from "@angular/core";
import { EventService } from "../services/event.service";
import { MatDialog } from "@angular/material";
import { SuccessDialogComponent } from "../dialogs/success/success-dialog.component";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.scss"]
})
export class AppComponent {
  constructor(private events: EventService, private dialog: MatDialog) {
    events.getEventListener().subscribe(event => {
      const data = JSON.parse(event.data);
      if (data.key === "Login") {
        if (data.value === "True") {
          const sDialogs = this.dialog.openDialogs.filter(d => {
            return d.componentInstance instanceof SuccessDialogComponent;
          });

          if (sDialogs.length === 0) {
            this.dialog.open(SuccessDialogComponent, {
              width: "75vw",
              data: {
                msg:
                  data.data + " (" + data.user + ") successfully checked in1!"
              }
            });
          }
          console.log(data);
        } else {
          // TODO: Show the failure dialog for a bad login
        }
      }
    });
  }
}
