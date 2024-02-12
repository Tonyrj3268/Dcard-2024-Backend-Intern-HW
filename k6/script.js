// import http from "k6/http";
// import { check, sleep } from "k6";
// import {
//   randomIntBetween,
//   randomItem,
// } from "https://jslib.k6.io/k6-utils/1.1.0/index.js";

// export let options = {
//   stages: [{ duration: "1s", target: 30000 }],
//   thresholds: {
//     http_req_duration: ["p(95)<1000"],
//   },
// };
// function randomBoolean() {
//   return Math.random() < 0.5;
// }
// function generateQueryParams() {
//   let queryParams = [];
//   // queryParams.push(`offset=${randomIntBetween(0, 100)}`);
//   // queryParams.push(`limit=${randomIntBetween(1, 100)}`);
//   queryParams.push(`offset=0`);
//   queryParams.push(`limit=10`);

//   // 條件性地添加其他參數
//   if (randomBoolean()) queryParams.push(`age=${randomIntBetween(1, 100)}`);
//   if (randomBoolean()) queryParams.push(`gender=${randomItem(["M", "F"])}`);
//   if (randomBoolean())
//     queryParams.push(`country=${randomItem(["US", "CN", "JP", "DE", "FR"])}`);
//   if (randomBoolean())
//     queryParams.push(`platform=${randomItem(["android", "ios", "web"])}`);

//   return queryParams.join("&");
// }

// export default function () {
//   let params = {
//     headers: {
//       "Content-Type": "application/json",
//     },
//   };

//   let queryParams = generateQueryParams();
//   let res = http.get(`http://nginx/api/v1/ad?${queryParams}`, params);
//   // let res = http.get(
//   //   `http://nginx/api/v1/ad?offset=0&limit=10&age=25&gender=M&country=US&platform=android`
//   // );

//   check(res, {
//     "is status 200": (r) => r.status === 200,
//   });

//   sleep(1);
// }

import {
  randomIntBetween,
  randomItem,
} from "https://jslib.k6.io/k6-utils/1.2.0/index.js";
import http from "k6/http";
import { sleep } from "k6";

export const options = {
  scenarios: {
    constant_request_rate: {
      executor: "constant-arrival-rate",
      rate: 24000,
      timeUnit: "1s",
      duration: "10s",
      preAllocatedVUs: 50,
      maxVUs: 10000,
    },
  },
};

let urlString = "http://" + "localhost" + "/api/v1/ad";

const limits = [5, 10, 15];
const genders = ["M", "F"];
const countries = ["TW", "JP"];
const platforms = ["android", "ios", "web"];

export default function () {
  const limit = randomItem(limits);
  const age = 25;
  const gender = randomItem(genders);
  const country = randomItem(countries);
  const platform = randomItem(platforms);

  for (let i = 0; i < 10; i++) {
    // 用 new URL 效能會變差
    http.get(
      `${urlString}?limit=${limit}&offset=${i}&age=${age}&gender=${gender}&country=${country}&platform=${platform}`
    );
  }

  sleep(1);
}
