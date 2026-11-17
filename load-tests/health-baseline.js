import http from 'k6/http';
import { check } from 'k6';

export const options = {
  vus: 1,
  duration: '30s',
};

export default function () {
  const response = http.get(`${__ENV.BASE_URL}/healthz`, {
    timeout: '5s',
  });

  check(response, {
    'status is 200': (r) => r.status === 200,
    'instance is correct': (r) => r.json('instance') === 'dsl-node-01',
    'status is ok': (r) => r.json('status') === 'ok',
  });
}
