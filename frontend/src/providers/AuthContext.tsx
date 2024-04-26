import {
  type FC,
  createContext,
  type ReactNode,
  useState,
  useEffect,
} from "react";
import useGetUser from "@/hooks/useGetUser";
import type User from "@/types/User";

interface AuthContextType {
  isAuthenticated: boolean;
  isLoading: boolean;
  user?: User;
}

export const AuthContext = createContext<AuthContextType>({
  isAuthenticated: false,
  isLoading: true,
  user: undefined,
});

const AuthProvider: FC<{ children: ReactNode }> = ({ children }) => {
  const [isAuthenticated, setAuthenticated] = useState<boolean>(false);
  const [user, setUser] = useState<User | undefined>(undefined);

  const { data, isLoading, isSuccess } = useGetUser();

  useEffect(() => {
    if (!isSuccess || !data) {
      setAuthenticated(false);
      return;
    }

    setAuthenticated(true);
    setUser(data);
  }, [data, isSuccess]);

  return (
    <AuthContext.Provider value={{ isAuthenticated, isLoading, user }}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthProvider;
