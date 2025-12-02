
export interface UserProfile {
  id: string;
  userId: string;
  name: string;
  firstName: string;
  lastName: string;
  dob: string;
  phone_number: string;
  gender: string;
  image?: {
    valid: boolean;
    data: string;
  };
}