package service

import (
	"github.com/jtejido/ngac/internal/set"
	// "github.com/jtejido/ngac/epp"
	// "github.com/jtejido/ngac/epp/events"
	// "github.com/jtejido/ngac/operations"
)

type AnalyticsService struct {
	Service
}

func NewAnalyticsService() *AnalyticsService {
	as := new(AnalyticsService)
	return as
}

func (as *AnalyticsService) Permissions(ctx Context, target string) set.Set {
	decider := as.Decider()
	return decider.List(ctx.User(), ctx.Process(), target)
}

// func (as *AnalyticsService) Pos(ctx context.Context) ([]*nodes.Node, error) {
//         // Prepare the hashset to return.
//         hsOa := make([]*nodes.Node, 0)

//         // Call find_border_oa_priv(u). The result is a Hashtable
//         // htoa = {oa -> {op -> pcset}}:
//          htOa := findBorderOaPrivRestrictedInternal(ctx);

//         // For each returned oa (key in htOa)
//         for (Enumeration oas = htOa.keys(); oas.hasMoreElements(); ) {
//             Node oa = (Node)oas.nextElement();

//             // Compute oa's required PCs by calling find_pc_set(sOaID).
//             HashSet<Long> hsReqPcs = inMemFindPcSet(oa.getID());
//             // Extract oa's label.
//             Hashtable htOaLabel = (Hashtable)htOa.get(oa);

//             // Walk through the op -> pcset of the oa's label.
//             // For each operation/access right
//             for (Enumeration ops = htOaLabel.keys(); ops.hasMoreElements(); ) {
//                 String sOp = (String)ops.nextElement();
//                 // Extract the pcset corresponding to this operation/access right.
//                 HashSet hsActualPcs = (HashSet)htOaLabel.get(sOp);
//                 // if the set of required PCs is a subset of the actual pcset,
//                 // then user u has some privileges on the current oa node.
//                 if (hsActualPcs.containsAll(hsReqPcs)) {
//                     hsOa.add(oa);
//                     break;
//                 }
//             }
//         }

//         return new HashSet<>(hsOa);
//     }

//     private Hashtable findBorderOaPrivRestrictedInternal( ctx context.Context) throws PMException {
//         // Uses a hashtable htReachableOas of reachable oas (see find_border_oa_priv(u))
//         // An oa is a key in this hashtable. The value is another hashtable that
//         // represents a label of the oa. A label is a set of pairs {(op -> pcset)}, with
//         // the op being the key and pcset being the value.
//         // {oa -> {op -> pcset}}.
//         Hashtable htReachableOas = new Hashtable();

//         // BFS from u (the base node). Prepare a queue.
//         Set<Long> visited = new HashSet<>();
//         Long crtNode;

//         // Get u's directly assigned attributes and put them into the queue.
//         Set<Long> hsAttrs = getGraphPAP().getParents(userCtx.getUserID());
//         List<Long> queue = new ArrayList<>(hsAttrs);

//         // While the queue has elements, extract an element from the queue
//         // and visit it.
//         while (!queue.isEmpty()) {
//             // Extract an ua from queue.
//             crtNode = queue.remove(0);
//             if (!visited.contains(crtNode)) {
//                 // If the ua has ua -> oa edges
//                 if (inMemUattrHasOpsets(crtNode)) {
//                     // Find the set of PCs reachable from ua.
//                     HashSet<Long> hsUaPcs = inMemFindPcSet(crtNode);

//                     // From each discovered ua traverse the edges ua -> oa.

//                     // Find the opsets of this user attribute. Note that the set of containers for this
//                     // node (user attribute) may contain not only opsets.
//                     Map<Long, OperationSet> assocs = getGraphPAP().getSourceAssociations(crtNode);

//                     // Go through the containers and only for opsets do the following.
//                     // For each opset ops of ua:
//                     for (Long targetID : assocs.keySet()) {
//                         // If oa is in htReachableOas
//                         if (htReachableOas.containsKey(targetID)) {
//                             // Then oa has a label op1 -> hsPcs1, op2 -> hsPcs2,...
//                             // Extract its label:
//                             Hashtable htOaLabel = (Hashtable)htReachableOas.get(targetID);

//                             // Get the operations from the opset:
//                             Set opers = assocs.get(targetID);
//                             // For each operation in the opset
//                             Iterator opersIter = opers.iterator();
//                             while (opersIter.hasNext()) {
//                                 String sOp = (String)opersIter.next();
//                                 // If the oa's label already contains the operation sOp
//                                 if (htOaLabel.containsKey(sOp)) {
//                                     // The label contains op -> some pcset.
//                                     // Do the union of the old pc with ua's pcset
//                                     HashSet hsPcs = (HashSet)htOaLabel.get(sOp);
//                                     hsPcs.addAll(hsUaPcs);
//                                 } else {
//                                     // The op is not in the oa's label.
//                                     // Create new op -> ua's pcs mappiing in the label.
//                                     HashSet hsNewPcs = new HashSet(hsUaPcs);
//                                     htOaLabel.put(sOp, hsNewPcs);
//                                 }
//                             }
//                         } else {
//                             // oa is not in htReachableOas.
//                             // Prepare a new label
//                             Hashtable htOaLabel = new Hashtable();

//                             // Get the operations from the opset:
//                             Set opers = assocs.get(targetID);
//                             // For each operation in the opset
//                             Iterator opersIter = opers.iterator();
//                             while (opersIter.hasNext()) {
//                                 String sOp = (String)opersIter.next();
//                                 // Add op -> pcs to the label.
//                                 HashSet hsNewPcs = new HashSet(hsUaPcs);
//                                 htOaLabel.put(sOp,  hsNewPcs);
//                             }

//                             // Add oa -> {op -> pcs}
//                             htReachableOas.put(targetID,  htOaLabel);
//                         }
//                     }
//                 }
//                 visited.add(crtNode);

//                 Set<Long> hsDescs = getGraphPAP().getParents(crtNode);
//                 queue.addAll(hsDescs);
//             }
//         }

//         // For each reachable oa in htReachableOas.keys
//         for (Enumeration keys = htReachableOas.keys(); keys.hasMoreElements() ;) {
//             Node oa = (Node)keys.nextElement();
//             // Compute {pc | oa ->+ pc}
//             HashSet hsOaPcs = inMemFindPcSet(oa.getID());
//             // Extract oa's label.
//             Hashtable htOaLabel = (Hashtable)htReachableOas.get(oa);
//             // The label contains op1 -> pcs1, op2 -> pcs2,...
//             // For each operation in the label
//             for (Enumeration lbl = htOaLabel.keys(); lbl.hasMoreElements();) {
//                 String sOp = (String)lbl.nextElement();
//                 // Intersect the pcset corresponding to this operation,
//                 // which comes from the uas, with the oa's pcset.
//                 HashSet oaPcs = (HashSet)htOaLabel.get(sOp);
//                 oaPcs.retainAll(hsOaPcs);
//                 if (oaPcs.isEmpty()) htOaLabel.remove(sOp);
//             }
//         }

//         return htReachableOas;
//     }

//     private HashSet<Long> inMemFindPcSet(Long node) throws PMException {
//         HashSet<Long> reachable = new HashSet<>();

//         // Init the queue, visited
//         ArrayList<Long> queue = new ArrayList<>();
//         HashSet<Long> visited = new HashSet<>();

//         // The current element
//         Long crtNode;

//         // Insert the start node into the queue
//         queue.add(node);

//         // While queue is not empty
//         while (!queue.isEmpty()) {
//             // Extract current element from queue
//             crtNode = queue.remove(0);
//             // If not visited
//             if (!visited.contains(crtNode)) {
//                 // Mark it as visited
//                 visited.add(crtNode);
//                 // Extract its direct descendants. If a descendant is an attribute,
//                 // insert it into the queue. If it is a pc, add it to reachable,
//                 // if not already there
//                 Set<Long> hsContainers = getGraphPAP().getParents(crtNode);
//                 Iterator<Long> hsiter = hsContainers.iterator();
//                 while (hsiter.hasNext()) {
//                     Long n = hsiter.next();
//                     if(getGraphPAP().getPolicyClasses().contains(n)) {
//                         reachable.add(n);
//                     } else {
//                         queue.add(n);
//                     }
//                 }
//             }
//         }
//         return reachable;
//     }

//     private boolean inMemUattrHasOpsets(Long uaNode) throws PMException {
//         return !getGraphPAP().getSourceAssociations(uaNode).isEmpty();
//     }

//     public Explain explain(long userID, long targetID) throws PMException {
//         Auditor auditor = new PReviewAuditor(getGraphPAP());
//         return auditor.explain(userID, targetID);
//     }
