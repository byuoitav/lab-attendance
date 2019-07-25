import { ByuIDPipe } from "./byu-id.pipe";

describe("ByuIdPipe", () => {
  let pipe: ByuIDPipe;

  beforeEach(() => {
    pipe = new ByuIDPipe();
  });

  it("create an instance", () => {
    expect(pipe).toBeTruthy();
  });

  it("should not modify a string less than 3 characters long", () => {
    expect(pipe.transform("")).toBe("");
    expect(pipe.transform("1")).toBe("1");
    expect(pipe.transform("12")).toBe("12");
  });

  it("should add only one hyphen for strings between 3 and 5 characters long", () => {
    expect(pipe.transform("123")).toBe("12-3");
    expect(pipe.transform("1234")).toBe("12-34");
    expect(pipe.transform("12345")).toBe("12-345");
  });

  it("should add two hyphens for strings larger than 6 characters", () => {
    expect(pipe.transform("123456")).toBe("12-345-6");
    expect(pipe.transform("1234567")).toBe("12-345-67");
    expect(pipe.transform("12345678")).toBe("12-345-678");
    expect(pipe.transform("123456789")).toBe("12-345-6789");
  });
});
