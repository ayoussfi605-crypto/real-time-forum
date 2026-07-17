import { setCurrentUser } from "./state.js";


export async function checkSession(){

    try{

        const response = await fetch("/me", { credentials: "include" });
        
        if (!response.ok){
            
            return false
        }
        
        const data = await response.json()

        setCurrentUser(data);
        return true;

    }catch{
        console.error(err);
        return false;
    }

}