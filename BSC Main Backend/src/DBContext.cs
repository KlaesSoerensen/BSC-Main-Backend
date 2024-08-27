using Microsoft.EntityFrameworkCore;

namespace BSC_Main_Backend.Models
{
    public class DBContext : DbContext
    {
        public DBContext(DbContextOptions<DBContext> options)
            : base(options)
        {
        }

        // Define a DbSet for each entity/table in your database
        public DbSet<GraphicalAsset> GraphicalAssets { get; set; }
        public DbSet<NPC> NPCs { get; set; }
        public DbSet<Transform> Transforms { get; set; }
        public DbSet<AssetCollection> AssetCollections { get; set; }
        public DbSet<Achievement> Achievements { get; set; }
        public DbSet<Player> Players { get; set; }
        public DbSet<Session> Sessions { get; set; }
        public DbSet<Colony> Colonies { get; set; }
        public DbSet<ColonyCode> ColonyCodes { get; set; }
        public DbSet<MiniGame> MiniGames { get; set; }
        public DbSet<MiniGameDifficulty> MiniGameDifficulties { get; set; }
        public DbSet<Location> Locations { get; set; }
        public DbSet<ColonyLocation> ColonyLocations { get; set; }
        public DbSet<ColonyAsset> ColonyAssets { get; set; }
        public DbSet<CollectionEntry> CollectionEntries { get; set; }
        public DbSet<LOD> LODs { get; set; }

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            base.OnModelCreating(modelBuilder);

            // Map the entities to the existing database tables
            modelBuilder.Entity<GraphicalAsset>().ToTable("GraphicalAsset");
            modelBuilder.Entity<NPC>().ToTable("NPC");
            modelBuilder.Entity<Transform>().ToTable("Transform");
            modelBuilder.Entity<AssetCollection>().ToTable("AssetCollection");
            modelBuilder.Entity<Achievement>().ToTable("Achievement");
            modelBuilder.Entity<Player>().ToTable("Player");
            modelBuilder.Entity<Session>().ToTable("Session");
            modelBuilder.Entity<Colony>().ToTable("Colony");
            modelBuilder.Entity<ColonyCode>().ToTable("ColonyCode");
            modelBuilder.Entity<MiniGame>().ToTable("MiniGame");
            modelBuilder.Entity<MiniGameDifficulty>().ToTable("MiniGameDifficulty");
            modelBuilder.Entity<Location>().ToTable("Location");
            modelBuilder.Entity<ColonyLocation>().ToTable("ColonyLocation");
            modelBuilder.Entity<ColonyAsset>().ToTable("ColonyAsset");
            modelBuilder.Entity<CollectionEntry>().ToTable("CollectionEntry");
            modelBuilder.Entity<LOD>().ToTable("LOD");

            // Relationships:
            modelBuilder.Entity<NPC>()
                .HasOne(n => n.SpriteNavigation)
                .WithMany()
                .HasForeignKey(n => n.Sprite)
                .OnDelete(DeleteBehavior.Restrict);
            // Each NPC has one Sprite (a GraphicalAsset), and the Sprite cannot be deleted if it's in use by an NPC.

            modelBuilder.Entity<Player>()
                .HasOne(p => p.SpriteNavigation)
                .WithMany()
                .HasForeignKey(p => p.Sprite)
                .OnDelete(DeleteBehavior.Restrict);
            // Each Player has one Sprite (a GraphicalAsset), and the Sprite cannot be deleted if it's in use by a Player.

            modelBuilder.Entity<Achievement>()
                .HasOne(a => a.IconNavigation)
                .WithMany()
                .HasForeignKey(a => a.Icon)
                .OnDelete(DeleteBehavior.Restrict);
            // Each Achievement is linked to one Icon (a GraphicalAsset), and the Icon cannot be deleted if it's in use by an Achievement.

            modelBuilder.Entity<Session>()
                .HasOne(s => s.PlayerNavigation)
                .WithMany()
                .HasForeignKey(s => s.Player)
                .OnDelete(DeleteBehavior.Cascade);
            // Each Session is associated with a Player, and deleting a Player will also delete their associated Sessions.

            modelBuilder.Entity<Colony>()
                .HasOne(c => c.OwnerNavigation)
                .WithMany()
                .HasForeignKey(c => c.Owner)
                .OnDelete(DeleteBehavior.Cascade);
            // Each Colony is owned by a Player, and deleting the Player will also delete their associated Colonies.

            modelBuilder.Entity<Colony>()
                .HasOne(c => c.ColonyCodeNavigation)
                .WithMany()
                .HasForeignKey(c => c.ColonyCode)
                .OnDelete(DeleteBehavior.SetNull);
            // Each Colony may have a ColonyCode, and deleting the ColonyCode will set the foreign key in the Colony to NULL.

            modelBuilder.Entity<ColonyCode>()
                .HasOne(cc => cc.ColonyNavigation)
                .WithMany()
                .HasForeignKey(cc => cc.Colony)
                .OnDelete(DeleteBehavior.Cascade);
            // Each ColonyCode is associated with a Colony, and deleting a Colony will delete its associated ColonyCode.

            modelBuilder.Entity<MiniGame>()
                .HasOne(mg => mg.IconNavigation)
                .WithMany()
                .HasForeignKey(mg => mg.Icon)
                .OnDelete(DeleteBehavior.Restrict);
            // Each MiniGame is linked to an Icon (a GraphicalAsset), and the Icon cannot be deleted if it's in use by a MiniGame.

            modelBuilder.Entity<MiniGameDifficulty>()
                .HasOne(mgd => mgd.MiniGameNavigation)
                .WithMany()
                .HasForeignKey(mgd => mgd.MiniGame)
                .OnDelete(DeleteBehavior.Cascade);
            // Each MiniGameDifficulty is associated with a MiniGame, and deleting a MiniGame will also delete its associated difficulties.

            modelBuilder.Entity<MiniGameDifficulty>()
                .HasOne(mgd => mgd.IconNavigation)
                .WithMany()
                .HasForeignKey(mgd => mgd.Icon)
                .OnDelete(DeleteBehavior.Restrict);
            // Each MiniGameDifficulty is linked to an Icon (a GraphicalAsset), and the Icon cannot be deleted if it's in use by a MiniGameDifficulty.

            modelBuilder.Entity<Location>()
                .HasOne(l => l.MiniGameNavigation)
                .WithMany()
                .HasForeignKey(l => l.MiniGame)
                .OnDelete(DeleteBehavior.SetNull);
            // Each Location may be associated with a MiniGame, and deleting the MiniGame will set the foreign key in the Location to NULL.

            modelBuilder.Entity<ColonyLocation>()
                .HasOne(cl => cl.ColonyNavigation)
                .WithMany()
                .HasForeignKey(cl => cl.Colony)
                .OnDelete(DeleteBehavior.Cascade);
            // Each ColonyLocation is associated with a Colony, and deleting the Colony will also delete its associated ColonyLocations.

            modelBuilder.Entity<ColonyLocation>()
                .HasOne(cl => cl.LocationNavigation)
                .WithMany()
                .HasForeignKey(cl => cl.Location)
                .OnDelete(DeleteBehavior.Cascade);
            // Each ColonyLocation is associated with a Location, and deleting the Location will also delete its associated ColonyLocations.

            modelBuilder.Entity<ColonyLocation>()
                .HasOne(cl => cl.TransformNavigation)
                .WithMany()
                .HasForeignKey(cl => cl.Transform)
                .OnDelete(DeleteBehavior.Cascade);
            // Each ColonyLocation is associated with a Transform, and deleting the Transform will also delete its associated ColonyLocations.

            modelBuilder.Entity<ColonyAsset>()
                .HasOne(ca => ca.AssetCollectionNavigation)
                .WithMany()
                .HasForeignKey(ca => ca.AssetCollection)
                .OnDelete(DeleteBehavior.Cascade);
            // Each ColonyAsset is associated with an AssetCollection, and deleting the AssetCollection will also delete its associated ColonyAssets.

            modelBuilder.Entity<ColonyAsset>()
                .HasOne(ca => ca.TransformNavigation)
                .WithMany()
                .HasForeignKey(ca => ca.Transform)
                .OnDelete(DeleteBehavior.Cascade);
            // Each ColonyAsset is associated with a Transform, and deleting the Transform will also delete its associated ColonyAssets.

            modelBuilder.Entity<ColonyAsset>()
                .HasOne(ca => ca.ColonyNavigation)
                .WithMany()
                .HasForeignKey(ca => ca.Colony)
                .OnDelete(DeleteBehavior.Cascade);
            // Each ColonyAsset is associated with a Colony, and deleting the Colony will also delete its associated ColonyAssets.

            modelBuilder.Entity<CollectionEntry>()
                .HasOne(ce => ce.TransformNavigation)
                .WithMany()
                .HasForeignKey(ce => ce.Transform)
                .OnDelete(DeleteBehavior.Cascade);
            // Each CollectionEntry is associated with a Transform, and deleting the Transform will also delete its associated CollectionEntries.

            modelBuilder.Entity<CollectionEntry>()
                .HasOne(ce => ce.AssetCollectionNavigation)
                .WithMany()
                .HasForeignKey(ce => ce.AssetCollection)
                .OnDelete(DeleteBehavior.Cascade);
            // Each CollectionEntry is associated with an AssetCollection, and deleting the AssetCollection will also delete its associated CollectionEntries.

            modelBuilder.Entity<CollectionEntry>()
                .HasOne(ce => ce.GraphicalAssetNavigation)
                .WithMany()
                .HasForeignKey(ce => ce.GraphicalAsset)
                .OnDelete(DeleteBehavior.Cascade);
            // Each CollectionEntry is associated with a GraphicalAsset, and deleting the GraphicalAsset will also delete its associated CollectionEntries.

            modelBuilder.Entity<LOD>()
                .HasOne(l => l.GraphicalAssetNavigation)
                .WithMany()
                .HasForeignKey(l => l.GraphicalAsset)
                .OnDelete(DeleteBehavior.Cascade);
            // Each LOD is associated with a GraphicalAsset, and deleting the GraphicalAsset will also delete its associated LODs.
        }
    }
}
