import { updateNavbar } from "./navbar.js";
import { setCurrentUser } from "./state.js";


export async function checkSession(){

    try{

        const response = await fetch("/me", { credentials: "include" });
        
        if (!response.ok){
            
            return false
        }
        
        const data = await response.json()
        
        setCurrentUser(data);
        updateNavbar()
        return true;

    }catch(err){
        console.error(err);
        return false;
    }

}