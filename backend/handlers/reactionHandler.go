package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    db "forum/database"
    "forum/helpers"
    "forum/middlewares"
)

type ReactionRequest struct {
    Reaction string `json:"reaction"`
}

func HandleReaction(w http.ResponseWriter, r *http.Request) {

    // 1. Read body
    var req ReactionRequest

    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        helpers.SendJSON(
            w,
            http.StatusBadRequest,
            "Invalid request body",
        )
        return
    }

    // 2. Validate reaction
    if req.Reaction != "like" && req.Reaction != "dislike" {
        helpers.SendJSON(
            w,
            http.StatusBadRequest,
            "Reaction must be like or dislike",
        )
        return
    }

    // 3. Get authenticated user
    userID := r.Context().Value(
        middlewares.UserIDKey,
    ).(int)

    // 4. Get post ID
    postID, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        helpers.SendJSON(
            w,
            http.StatusBadRequest,
            "Invalid post ID",
        )
        return
    }

    // 5. Check existing reaction
    reaction, err := db.GetReaction(postID, userID)
    if err != nil {
        helpers.SendJSON(
            w,
            http.StatusInternalServerError,
            "Failed to get reaction",
        )
        return
    }

    userReaction := req.Reaction
    message := ""

    // 6. No reaction → CREATE
    if reaction == nil {

        err := db.CreateReaction(
            postID,
            userID,
            req.Reaction,
        )

        if err != nil {
            helpers.SendJSON(
                w,
                http.StatusInternalServerError,
                "Failed to create reaction",
            )
            return
        }

        message = "Reaction created"

    } else {

        // 7. Same reaction → DELETE
        if reaction.Reaction == req.Reaction {

            err := db.DeleteReaction(
                postID,
                userID,
            )

            if err != nil {
                helpers.SendJSON(
                    w,
                    http.StatusInternalServerError,
                    "Failed to delete reaction",
                )
                return
            }

            userReaction = ""
            message = "Reaction deleted"

        } else {

            // 8. Different reaction → UPDATE
            err := db.UpdateReaction(
                postID,
                userID,
                req.Reaction,
            )

            if err != nil {
                helpers.SendJSON(
                    w,
                    http.StatusInternalServerError,
                    "Failed to update reaction",
                )
                return
            }

            message = "Reaction updated"
        }
    }

    // 9. Count reactions
    stats, err := db.CountReactions(postID)

    if err != nil {
        helpers.SendJSON(
            w,
            http.StatusInternalServerError,
            "Failed to count reactions",
        )
        return
    }

    // 10. Response
    helpers.SendJSON(
        w,
        http.StatusOK,
        message,
        map[string]interface{}{
            "likes":        stats.Likes,
            "dislikes":     stats.Dislikes,
            "userReaction": userReaction,
        },
    )
}