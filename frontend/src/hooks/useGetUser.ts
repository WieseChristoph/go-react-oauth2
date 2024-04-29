import { useQuery } from "react-query";
import type User from "@/types/User";
import apiRequest from "@/utils/apiRequest";

export const GET_USER_QUERY_KEY = "getUser";

async function fetchUser() {
  return apiRequest("/users/me", "GET");
}

function useGetUser() {
  return useQuery<User, Error>(GET_USER_QUERY_KEY, fetchUser, {
    refetchOnWindowFocus: false,
  });
}

export default useGetUser;
