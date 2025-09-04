import axios from 'axios';

export async function fetchStatus(url: string): Promise<any> {
  const response = await axios.get(url);
  return response.data;
}
