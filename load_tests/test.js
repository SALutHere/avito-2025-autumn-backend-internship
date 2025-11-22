import http from "k6/http";
import { sleep, check } from "k6";

export let options = {
    stages: [
        { duration: "30s", target: 200 },
        { duration: "1m",  target: 200 },
        { duration: "30s", target: 0 },
    ],
};

export default function () {
    const res = http.get("http://app:8080/stats");
    check(res, {
        "status 200": (r) => r.status === 200,
        "time < 200ms": (r) => r.timings.duration < 200,
    });

    sleep(1);
}
