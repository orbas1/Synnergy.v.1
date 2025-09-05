import { renderHome } from './pages/Home';

export function main(balance = 0): string {
  return renderHome(balance);
}

if (require.main === module) {
  console.log(main());
}
