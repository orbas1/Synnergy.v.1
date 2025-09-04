export class IdentityService {
  private users = new Map<string, string>();

  register(user: string, publicKey: string): void {
    if (this.users.has(user)) {
      throw new Error('user already registered');
    }
    this.users.set(user, publicKey);
  }

  getUser(user: string): string | undefined {
    return this.users.get(user);
  }
}
