import {useEffect} from "react";
import {useDispatch, useSelector} from "react-redux";
import {getBots} from "../../features/bots/botsSlice";
import BotCard from "./BotCard";

export default function BotsList()
{
    const dispatch=useDispatch();
    const {list, loading}=useSelector(state => state.bots);

    useEffect(() =>
    {
        dispatch(getBots());
    }, []);

    if (loading) return <p>Loading...</p>;

    return (
        <div>
            {list.map(bot => (
                <BotCard key={bot.id} bot={bot} />
            ))}
        </div>
    );
}
