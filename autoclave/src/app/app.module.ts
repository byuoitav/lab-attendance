import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";
import { HttpClientModule } from "@angular/common/http";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import {
  MatToolbarModule,
  MatButtonModule,
  MatIconModule,
  MatDialogModule
} from "@angular/material";

import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./components/app.component";
import { LoginComponent } from "./components/login/login.component";
import { ScreenSaverComponent } from "./components/screen-saver/screen-saver.component";
import { ByuIDPipe } from "./pipes/byu-id.pipe";
import { SuccessDialogComponent } from "./dialogs/success/success-dialog.component";

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    ScreenSaverComponent,
    ByuIDPipe,
    SuccessDialogComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    MatToolbarModule,
    MatButtonModule,
    MatIconModule,
    MatDialogModule,
    HttpClientModule
  ],
  providers: [],
  entryComponents: [SuccessDialogComponent],
  bootstrap: [AppComponent]
})
export class AppModule {}
