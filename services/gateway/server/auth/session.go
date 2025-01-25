package authhandler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/common"
	"github.com/wafi04/golang-backend/services/common/middleware"
)


func (s *AuthHandler) HandlerListSessions(w http.ResponseWriter, r *http.Request) {
    user, err :=  middleware.GetUserFromContext(r.Context())
    if err != nil {
        common.Error(http.StatusUnauthorized, "Unauthorized")
        return
    }

    listSessions, err := s.authClient.ListSessions(r.Context(), &pb.ListSessionsRequest{
        UserId: user.UserId,
    })
    if err != nil {
        common.Error(http.StatusInternalServerError, "Failed to retrieve sessions")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    
    res := common.Success(listSessions, "Sessions Retrieved Successfully")
    if err := json.NewEncoder(w).Encode(res); err != nil {
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
        return
    }
}



func (h *AuthHandler)  HandleRevokeSessions(w http.ResponseWriter, r *http.Request){
    user,err :=  middleware.GetUserFromContext(r.Context())

    vars := mux.Vars(r)
    session, ok := vars["id"]
    if !ok {
        http.Error(w, "Category ID is required", http.StatusBadRequest)
        return
    }

    if err != nil {
        common.Error(http.StatusUnauthorized, "Unauthorized")
    }

    revoke,err :=  h.authClient.RevokeSession(r.Context(), &pb.RevokeSessionRequest{
        UserId:user.UserId,
        SessionId: session,
    })


    if err != nil {
        common.Error(http.StatusUnauthorized, "Failed To Delete Session")
    }

    res :=  common.Success(revoke, "Delete Succes")

     if err := json.NewEncoder(w).Encode(res); err != nil {
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
        return
    }
}