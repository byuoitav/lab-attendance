import { Injectable } from "@angular/core";
import { HttpClient, HttpHeaders } from "@angular/common/http";

@Injectable({
  providedIn: "root"
})
export class APIService {
  constructor(private http: HttpClient) {}

  login = (byuID: string) => {
    const endpoint =
      "http://" + window.location.host + "/api/v1/login/" + byuID;

    console.log("Hitting endpoint: " + endpoint);

    this.http
      .post(endpoint, null, {})
      .subscribe(
        (data: object) => console.log(data),
        error => console.log(error)
      );
  };
}
